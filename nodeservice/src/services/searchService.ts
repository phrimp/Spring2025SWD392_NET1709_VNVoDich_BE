import { PrismaClient } from "@prisma/client";

const prisma = new PrismaClient();

export const findTutors = async (query: string, pageNum: number, pageSizeNum: number) => {
  const skip = (pageNum - 1) * pageSizeNum;
  return await prisma.tutor.findMany({
    where: { profile: { full_name: { contains: query } } },
    skip,
    take: pageSizeNum,
    include: { profile: true },
  });
};

export const findCourses = async (query: string, pageNum: number, pageSizeNum: number) => {
  const skip = (pageNum - 1) * pageSizeNum;
  return await prisma.course.findMany({
    where: { title: { contains: query } },
    skip,
    take: pageSizeNum,
  });
};

export const filterCoursesByPrice = async (minPrice: number, maxPrice: number) => {
  return await prisma.course.findMany({
    where: {
      price: {
        gte: minPrice,
        lte: maxPrice,
      },
    },
  });
};

export const filterCoursesByGrade = async (grade: number) => {
  return await prisma.course.findMany({
    where: { grade: { equals: grade } },
  });
};

export const filterCoursesBySubject = async (subject: string) => {
  return await prisma.course.findMany({
    where: { subject: { contains: subject } },
  });
};

export const filterTutorsByRating = async (minRating: number) => {
  const tutorsWithRatings = await prisma.tutorReview.groupBy({
    by: ["tutor_id"],
    _avg: { rating: true },
    having: {
      rating: { _avg: { gte: minRating } }, // Lọc gia sư có rating trung bình >= minRating
    },
  });

  // Lấy danh sách Tutor ID từ kết quả groupBy
  const tutorIds = tutorsWithRatings.map((t) => t.tutor_id);

  // Truy vấn chi tiết thông tin gia sư
  const tutors = await prisma.tutor.findMany({
    where: { id: { in: tutorIds } },
    include: {
      profile: {
        select: {
          full_name: true,
          email: true,
          phone: true,
        },
      },
    },
  });

  // Ghép thông tin rating trung bình vào kết quả
  return tutors.map((tutor) => {
    const avgRating =
      tutorsWithRatings.find((r) => r.tutor_id === tutor.id)?._avg?.rating ?? 0;
    return { ...tutor, avgRating };
  });
};
