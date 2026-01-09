let lastTitle = '';

export const useTitle = (title: string) => {
  lastTitle = document.title;
  document.title = `${title} | Vicar`;
};

export const resetTitle = () => {
  document.title = lastTitle;
};