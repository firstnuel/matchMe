export const calculateAge = (birthdate: string|number): number => {
  const birth = new Date(birthdate);
  const today = new Date();

  let age = today.getFullYear() - birth.getFullYear();

  // Adjust if birthday hasn't happened yet this year
  const hasHadBirthdayThisYear =
    today.getMonth() > birth.getMonth() ||
    (today.getMonth() === birth.getMonth() && today.getDate() >= birth.getDate());

  if (!hasHadBirthdayThisYear) {
    age--;
  }

  if (age < 1) return 1 // to avoid go negative and 0 error

  return age;
}

