export function interpolateString(str: string, obj: any): string {
  return str.replace(/{([^{}]*)}/g, (a: string, b: string) => {
    const r = obj[b];
    return typeof r === 'string' || typeof r === 'number' ? String(r) : b;
  }
  );
}