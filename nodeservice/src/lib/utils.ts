export const calculateAge = (dateOfBirth: string) => {
  const birthDate = new Date(dateOfBirth.split("T")[0]); // Lấy phần YYYY-MM-DD
  const today = new Date();

  let age = today.getFullYear() - birthDate.getFullYear();
  const monthDiff = today.getMonth() - birthDate.getMonth();
  const dayDiff = today.getDate() - birthDate.getDate();

  // Điều chỉnh tuổi nếu sinh nhật chưa tới trong năm hiện tại
  if (monthDiff < 0 || (monthDiff === 0 && dayDiff < 0)) {
    age--;
  }
  return age;
};
