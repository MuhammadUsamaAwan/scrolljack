import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function pascalCaseToTitleCase(str: string): string {
  return str
    .replace(/([a-z])([A-Z])/g, '$1 $2')
    .replace(/([A-Z])([A-Z][a-z])/g, '$1 $2')
    .replace(/([0-9]+)/g, ' $1 ')
    .trim()
    .replace(/\s+/g, ' ')
    .toLowerCase()
    .replace(/(^\w|\s\w)/g, m => m.toUpperCase());
}
