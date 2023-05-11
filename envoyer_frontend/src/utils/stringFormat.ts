export const transformToVariable = (text: string) => {
  if (text.startsWith('{{') && text.endsWith('}}')) return text;
  return `{{.${text}}}`;
};
