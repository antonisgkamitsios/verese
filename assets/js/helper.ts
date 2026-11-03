export function parseHTML(str: string) {
  const parser = new DOMParser();
  return parser.parseFromString(str, 'text/html');
}