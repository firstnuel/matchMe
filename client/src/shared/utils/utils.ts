export const firstToUpper = (s: string): string => {
    return s.length <= 1? s : s.charAt(0).toUpperCase() + s.substring(1, s.length)
}

export function getInitials(firstName: string, lastName: string): string {
  return (firstName.charAt(0) + lastName.charAt(0)).toUpperCase();
}