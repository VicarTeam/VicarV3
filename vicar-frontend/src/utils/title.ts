const currentTitle = window.document.title;

export function setTitle(title: string) {
  window.document.title = `${title} | Vicar`;
}

export function resetTitle() {
  window.document.title = currentTitle;
}