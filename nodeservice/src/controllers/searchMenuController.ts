import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";

const prisma = new PrismaClient();

// export const searchTutorsAndCourses = async (
//   req: Request,
//   res: Response
// ): Promise<void> => {
//   try {
//     const { query, page = 1, pageSize = 10 } = req.query;

//     const pageNum = parseInt(page as string, 10);
//     const pageSizeNum = parseInt(pageSize as string, 10);
//     const skip = (pageNum - 1) * pageSizeNum;

//     const tutors = await prisma.tutor.findMany({
//       where: { profile: { full_name: { contains: query as string } } },
//       skip,
//       take: pageSizeNum,
//       include: { profile: true },
//     });

//     const courses = await prisma.course.findMany({
//       where: { title: { contains: query as string } },
//       skip,
//       take: pageSizeNum,
//     });

//     res.json({
//       message: "Search results retrieved successfully",
//       data: { tutors, courses },
//     });
//   } catch (error) {
//     res.status(500).json({ message: "Error searching", error });
//   }
// };

export const searchTutors = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { query, page = 1, pageSize = 10 } = req.query;

    const pageNum = parseInt(page as string, 10);
    const pageSizeNum = parseInt(pageSize as string, 10);
    const skip = (pageNum - 1) * pageSizeNum;
    const tutors = await prisma.tutor.findMany({
      where: { profile: { full_name: { contains: query as string } } },
      skip,
      take: pageSizeNum,
      include: { profile: true },
    });
    res.json({
      message: "Search results retrieved successfully",
      data: { tutors },
    });
  } catch (error) {
    res.status(500).json({ message: "Error searching tutors", error });
  }
};
export const searchCourses = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { query, page = 1, pageSize = 10 } = req.query;
    const pageNum = parseInt(page as string, 10);
    const pageSizeNum = parseInt(pageSize as string, 10);
    const skip = (pageNum - 1) * pageSizeNum;
    const courses = await prisma.course.findMany({
      where: { title: { contains: query as string } },
      skip,
      take: pageSizeNum,
    });

    res.json({
      message: "Search results retrieved successfully",
      data: { courses },
    });
  } catch (error) {
    res.status(500).json({ message: "Error searching courses", error });
  }
};
// Tìm kiếm lịch trình có sẵn
// export const searchSchedules = async (
//   req: Request,
//   res: Response
// ): Promise<void> => {
//   try {
//     const { date } = req.query;

//     const schedules = await prisma.schedule.findMany({
//       where: { date: { equals: new Date(date as string) } },
//     });

//     res.json({ message: "Schedules retrieved successfully", data: schedules });
//   } catch (error) {
//     res.status(500).json({ message: "Error retrieving schedules", error });
//   }
// };

// Lọc theo giá
export const filterByPrice = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { minPrice, maxPrice } = req.query;

    const courses = await prisma.course.findMany({
      where: {
        price: {
          gte: Number(minPrice),
          lte: Number(maxPrice),
        },
      },
    });

    res.json({ message: "Courses filtered by price", data: courses });
  } catch (error) {
    res.status(500).json({ message: "Error filtering by price", error });
  }
};

// Tìm kiếm theo cấp học
export const searchByGrade = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { grade } = req.query;

    const courses = await prisma.course.findMany({
      where: { grade: { equals: Number(grade) } },
    });

    res.json({ message: "Courses filtered by grade", data: courses });
  } catch (error) {
    res.status(500).json({ message: "Error filtering by grade", error });
  }
};

// Tìm kiếm theo môn học
export const searchBySubject = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { subject } = req.query;

    const courses = await prisma.course.findMany({
      where: { subject: { contains: subject as string } },
    });

    res.json({ message: "Courses filtered by subject", data: courses });
  } catch (error) {
    res.status(500).json({ message: "Error filtering by subject", error });
  }
};

//Lọc gia sư theo đánh giá
export const filterTutorsByRating = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const minRating = Number(req.query.minRating) || 0;

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
    const tutorsWithAvgRating = tutors.map((tutor) => {
      const avgRating =
        tutorsWithRatings.find((r) => r.tutor_id === tutor.id)?._avg?.rating ??
        0;
      return { ...tutor, avgRating };
    });

    res.json({
      message: "Tutors filtered by rating",
      data: tutorsWithAvgRating,
    });
  } catch (error) {
    console.error("Error filtering tutors by rating:", error);
    res.status(500).json({ message: "Error filtering tutors", error });
  }
};

//hehehehe
