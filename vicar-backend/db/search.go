package db

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SearchBuilder struct {
	debug        bool
	model        interface{}
	itemsPerPage int
	page         int
	preloads     []string
	joins        *[]struct {
		clause string
		args   []interface{}
	}
	innerJoins *[]struct {
		clause string
		args   []interface{}
	}
	whereClause *[]struct {
		condition string
		args      []interface{}
	}
	searchFilters *[]struct {
		name         string
		exact        bool
		ignoreCase   bool
		valueIsArray bool
		overrideDef  bool
	}
	preFilters *[]struct {
		key   string
		where string
	}
	preFilterValues   map[string]string
	orderBy           map[string]string
	searchTerm        string
	unscoped          bool
	disablePagination bool
}

// Creates a new search builder with default values.
// The search builder is used to build search queries and paginate the results. The search filtering is optional. The builder can also be used only for pagination.
func NewSearchBuilder(model interface{}) *SearchBuilder {
	return &SearchBuilder{
		model:             model,
		debug:             false,
		itemsPerPage:      20,
		page:              1,
		searchTerm:        "",
		preloads:          []string{},
		orderBy:           map[string]string{},
		preFilterValues:   map[string]string{},
		unscoped:          false,
		disablePagination: false,
	}
}

func (s *SearchBuilder) Debug() *SearchBuilder {
	s.debug = true
	return s
}

func (s *SearchBuilder) IgnoreSoftDelete() *SearchBuilder {
	s.unscoped = true
	return s
}

func (s *SearchBuilder) DisablePagination() *SearchBuilder {
	s.disablePagination = true
	return s
}

// Adds a new pre filter to the search query. The pre filter is used to filter the results before the search filters are applied. The key is used to retrieve pre filter values
// from the query parameters.
func (s *SearchBuilder) PreFilter(key, where string) *SearchBuilder {
	if s.preFilters == nil {
		s.preFilters = &[]struct {
			key   string
			where string
		}{}
	}

	*s.preFilters = append(*s.preFilters, struct {
		key   string
		where string
	}{
		key:   key,
		where: where,
	})

	return s
}

// Adds a new search filter. The search filter is used to filter the results by a column name. The exact parameter is used to determine if the search should be exact or using a LIKE clause.
// The ignoreCase parameter is used to determine if the search should be case sensitive or not.
func (s *SearchBuilder) SearchFilter(name string, exact bool, ignoreCase bool) *SearchBuilder {
	if s.searchFilters == nil {
		s.searchFilters = &[]struct {
			name         string
			exact        bool
			ignoreCase   bool
			valueIsArray bool
			overrideDef  bool
		}{}
	}

	*s.searchFilters = append(*s.searchFilters, struct {
		name         string
		exact        bool
		ignoreCase   bool
		valueIsArray bool
		overrideDef  bool
	}{
		name:         name,
		exact:        exact,
		ignoreCase:   ignoreCase,
		valueIsArray: false,
	})

	return s
}

// Adds a new array search filter. The array search filter is used to filter the results by a column name that is an array. The exact parameter is used to determine if the search should be exact or using a LIKE clause.
// The ignoreCase parameter is used to determine if the search should be case sensitive or not.
// The array search filter is used to filter the results by a column name that is an array. In this case unnest is used to extract the array elements and filter them.
func (s *SearchBuilder) ArraySearchFilter(name string, exact bool, ignoreCase bool) *SearchBuilder {
	if s.searchFilters == nil {
		s.searchFilters = &[]struct {
			name         string
			exact        bool
			ignoreCase   bool
			valueIsArray bool
			overrideDef  bool
		}{}
	}

	*s.searchFilters = append(*s.searchFilters, struct {
		name         string
		exact        bool
		ignoreCase   bool
		valueIsArray bool
		overrideDef  bool
	}{
		name:         name,
		exact:        exact,
		ignoreCase:   ignoreCase,
		valueIsArray: true,
	})

	return s
}

func (s *SearchBuilder) OverrideDefaultSearchFilter(filter string) *SearchBuilder {
	if s.searchFilters == nil {
		s.searchFilters = &[]struct {
			name         string
			exact        bool
			ignoreCase   bool
			valueIsArray bool
			overrideDef  bool
		}{}
	}

	*s.searchFilters = append(*s.searchFilters, struct {
		name         string
		exact        bool
		ignoreCase   bool
		valueIsArray bool
		overrideDef  bool
	}{
		name:         filter,
		exact:        false,
		ignoreCase:   false,
		valueIsArray: false,
		overrideDef:  true,
	})

	return s
}

// Extracts the itemsPerPage and page parameters from the URL query parameters.
func (s *SearchBuilder) Extract(c *fiber.Ctx) *SearchBuilder {
	if c.Query("itemsPerPage") != "" {
		s.itemsPerPage, _ = strconv.Atoi(c.Query("itemsPerPage"))
	}

	if c.Query("page") != "" {
		s.page, _ = strconv.Atoi(c.Query("page"))
	}

	if c.Query("search") != "" {
		s.searchTerm = c.Query("search")
	}

	for key, value := range c.Queries() {
		if strings.HasPrefix(key, "pre[") && strings.HasSuffix(key, "]") {
			key = key[4 : len(key)-1]
			s.preFilterValues[key] = value
		}
	}

	return s
}

// Sets the number of items per page.
func (s *SearchBuilder) ItemsPerPage(itemsPerPage int) *SearchBuilder {
	s.itemsPerPage = itemsPerPage
	return s
}

// Sets the page number. The page number is 1-based.
func (s *SearchBuilder) Page(page int) *SearchBuilder {
	s.page = page
	return s
}

// Sets the search term. The search term is used to filter the results by a column name. The exact parameter is used to determine if the search should be exact or using a LIKE clause.
func (s *SearchBuilder) Search(searchTerm string) *SearchBuilder {
	s.searchTerm = searchTerm
	return s
}

// Sets the where clause. The condition is a string that can contain placeholders. The placeholders are replaced by the given arguments.
func (s *SearchBuilder) Where(condition string, args ...interface{}) *SearchBuilder {
	if s.whereClause == nil {
		s.whereClause = &[]struct {
			condition string
			args      []interface{}
		}{}
	}

	*s.whereClause = append(*s.whereClause, struct {
		condition string
		args      []interface{}
	}{
		condition: condition,
		args:      args,
	})
	return s
}

// Sets the where clause. The condition is a string that can contain placeholders. The placeholders are replaced by the given arguments.
// The conditions are concatenated using the given gate parameter. The gate parameter is used to determine the logic gate to use between the conditions.
// The parts parameter is a map where the key is the predicate and the value is a slice of values to use in the condition.
func (s *SearchBuilder) WhereConcat(conditionParts map[string][]interface{}, gate string) *SearchBuilder {
	parts := make([]string, len(conditionParts))
	args := []interface{}{}

	i := 0
	for predicate, values := range conditionParts {
		parts[i] = predicate
		args = append(args, values...)

		i++
	}

	return s.Where(strings.Join(parts, gate), args...)
}

// Sets the joins clause. The joins clause is a string that can contain placeholders. The placeholders are replaced by the given arguments.
func (s *SearchBuilder) Joins(joins string, args ...interface{}) *SearchBuilder {
	if s.joins == nil {
		s.joins = &[]struct {
			clause string
			args   []interface{}
		}{}
	}

	*s.joins = append(*s.joins, struct {
		clause string
		args   []interface{}
	}{
		clause: joins,
		args:   args,
	})
	return s
}

// Sets the inner joins clause. The inner joins clause is a string that can contain placeholders. The placeholders are replaced by the given arguments.
func (s *SearchBuilder) InnerJoins(innerJoins string, args ...interface{}) *SearchBuilder {
	if s.innerJoins == nil {
		s.innerJoins = &[]struct {
			clause string
			args   []interface{}
		}{}
	}

	*s.innerJoins = append(*s.innerJoins, struct {
		clause string
		args   []interface{}
	}{
		clause: innerJoins,
		args:   args,
	})
	return s
}

// Sets the preloads clause. The preloads clause is a string that can contain placeholders. The placeholders are replaced by the given arguments.
func (s *SearchBuilder) Preload(preloads ...string) *SearchBuilder {
	s.preloads = append(s.preloads, preloads...)
	return s
}

// Sets the order by clause. The column is the name of the column to order by. The direction is the direction of the ordering (ASC or DESC)
func (s *SearchBuilder) OrderBy(column string, direction string) *SearchBuilder {
	s.orderBy[column] = direction
	return s
}

func (s *SearchBuilder) OrderByAsc(column string) *SearchBuilder {
	s.orderBy[column] = "ASC"
	return s
}

func (s *SearchBuilder) OrderByDesc(column string) *SearchBuilder {
	s.orderBy[column] = "DESC"
	return s
}

// Executes the search query. The results are bound to the given pointer to a slice. The total number of possible results is returned.
// If an error occurs, the number of results is 0 and the error is returned. If no results are found, the number of results is 0 and the error is nil.
func (s *SearchBuilder) Execute(results interface{}) (int64, error) {
	rctx := DB.Model(s.model)
	cctx := DB.Model(s.model)

	if s.debug {
		rctx = rctx.Debug()
		cctx = cctx.Debug()
	}

	if s.unscoped {
		rctx = rctx.Unscoped()
		cctx = cctx.Unscoped()
	}

	for _, preload := range s.preloads {
		rctx = rctx.Preload(preload)
		cctx = cctx.Preload(preload)
	}

	if s.preFilters != nil {
		for _, filter := range *s.preFilters {
			val, ok := s.preFilterValues[filter.key]
			if !ok {
				continue
			}

			questionMarks := strings.Count(filter.where, "?")
			args := []interface{}{val}
			if questionMarks > 1 {
				for i := 1; i < questionMarks; i++ {
					args = append(args, val)
				}
			}

			rctx = rctx.Where(filter.where, args...)
			cctx = cctx.Where(filter.where, args...)
		}

	}

	if s.whereClause != nil {
		for _, clause := range *s.whereClause {
			rctx = rctx.Where(clause.condition, clause.args...)
			cctx = cctx.Where(clause.condition, clause.args...)
		}
	}

	if s.joins != nil {
		for _, clause := range *s.joins {
			rctx = rctx.Joins(clause.clause, clause.args...)
			cctx = cctx.Joins(clause.clause, clause.args...)
		}
	}

	if s.innerJoins != nil {
		for _, clause := range *s.innerJoins {
			rctx = rctx.Joins(clause.clause, clause.args...)
			cctx = cctx.Joins(clause.clause, clause.args...)
		}
	}

	var queryBuilder strings.Builder
	var queryParams []interface{}

	searchTerm := s.searchTerm
	searchPrefix := ""
	logicGate := " OR "
	isExact := false
	whenExactIgnoreCase := true
	if strings.HasPrefix(searchTerm, "-") {
		searchPrefix = "NOT "
		logicGate = " AND "
		searchTerm = searchTerm[1:]
	} else if strings.HasPrefix(searchTerm, "=") {
		isExact = true

		if strings.HasPrefix(searchTerm, "=!") {
			whenExactIgnoreCase = false
			searchTerm = searchTerm[2:]
		} else {
			searchTerm = searchTerm[1:]
		}
	}

	if s.searchFilters != nil && s.searchTerm != "" {
		for i, filter := range *s.searchFilters {
			if i > 0 {
				queryBuilder.WriteString(logicGate)
			}

			if filter.overrideDef {
				queryBuilder.WriteString(filter.name)

				count := strings.Count(filter.name, "?")
				if count > 1 {
					for i := 1; i <= count; i++ {
						queryParams = append(queryParams, "%"+searchTerm+"%")
					}
				} else {
					queryParams = append(queryParams, "%"+searchTerm+"%")
				}
			} else {
				if !filter.valueIsArray {
					if filter.exact || isExact {
						if filter.ignoreCase && (!isExact || whenExactIgnoreCase) {
							queryBuilder.WriteString(fmt.Sprintf("%s(LOWER(%s) = LOWER(?))", searchPrefix, filter.name))
							queryParams = append(queryParams, strings.ToLower(searchTerm))
						} else {
							queryBuilder.WriteString(fmt.Sprintf("%s(%s = ?)", searchPrefix, filter.name))
							queryParams = append(queryParams, searchTerm)
						}
					} else {
						if filter.ignoreCase {
							queryBuilder.WriteString(fmt.Sprintf("%s(LOWER(%s) LIKE LOWER(?))", searchPrefix, filter.name))
							queryParams = append(queryParams, "%"+strings.ToLower(searchTerm)+"%")
						} else {
							queryBuilder.WriteString(fmt.Sprintf("%s(%s LIKE ?)", searchPrefix, filter.name))
							queryParams = append(queryParams, "%"+searchTerm+"%")
						}
					}
				} else {
					if filter.exact || isExact {
						if filter.ignoreCase && (!isExact || whenExactIgnoreCase) {
							queryBuilder.WriteString(fmt.Sprintf("%sEXISTS (SELECT 1 FROM unnest(%s) AS element WHERE LOWER(element) = LOWER(?))", searchPrefix, filter.name))
							queryParams = append(queryParams, strings.ToLower(searchTerm))
						} else {
							queryBuilder.WriteString(fmt.Sprintf("%sEXISTS (SELECT 1 FROM unnest(%s) AS element WHERE element = ?)", searchPrefix, filter.name))
							queryParams = append(queryParams, searchTerm)
						}
					} else {
						if filter.ignoreCase {
							queryBuilder.WriteString(fmt.Sprintf("%sEXISTS (SELECT 1 FROM unnest(%s) AS element WHERE LOWER(element) LIKE LOWER(?))", searchPrefix, filter.name))
							queryParams = append(queryParams, "%"+strings.ToLower(searchTerm)+"%")
						} else {
							queryBuilder.WriteString(fmt.Sprintf("%sEXISTS (SELECT 1 FROM unnest(%s) AS element WHERE element LIKE ?)", searchPrefix, filter.name))
							queryParams = append(queryParams, "%"+searchTerm+"%")
						}
					}
				}
			}
		}

		rctx = rctx.Where(queryBuilder.String(), queryParams...)
		cctx = cctx.Where(queryBuilder.String(), queryParams...)
	}

	for column, direction := range s.orderBy {
		rctx = rctx.Order(column + " " + direction)
	}

	if !s.disablePagination {
		if s.itemsPerPage > 0 {
			rctx = rctx.Limit(s.itemsPerPage)
		}

		if s.itemsPerPage > 0 {
			rctx = rctx.Offset((s.page - 1) * s.itemsPerPage)
		}
	}

	if res := rctx.Find(results); res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return 0, nil
		}

		return 0, res.Error
	}

	var count int64

	if !s.disablePagination {
		if res := cctx.Count(&count); res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				return 0, nil
			}

			return 0, res.Error
		}
	} else {
		count = -1
	}

	return count, nil
}
