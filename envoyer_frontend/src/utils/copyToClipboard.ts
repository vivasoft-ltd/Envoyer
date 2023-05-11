import { toast } from 'react-toastify';

// export const copyToClipboard = (text: string) => {
//   navigator.clipboard.writeText(text);
//   toast.info('Copied to clipboard');
// };

export const copyToClipboard = (text: string) => {
  var input = document.createElement('textarea');
  input.innerHTML = text;
  document.body.appendChild(input);
  input.select();
  var result = document.execCommand('copy');
  document.body.removeChild(input);
  toast.info('Copied to clipboard');

  return result;
};
