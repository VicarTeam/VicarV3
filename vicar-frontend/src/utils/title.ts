const currentTitle = window.document.title;

export function setTitle(title: string) {
  window.document.title = `Nauri - ${title}`;
}

export function resetTitle() {
  window.document.title = currentTitle;
}