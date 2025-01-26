import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";

const prisma = new PrismaClient();

export const getCourses = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    // Lấy query từ request
    const { page = 1, pageSize = 10, subject, grade, status } = req.query;

    // Chuyển đổi page và pageSize sang số nguyên
    const pageNum = parseInt(page as string, 10);
    const pageSizeNum = parseInt(pageSize as string, 10);

    // Tạo bộ lọc dựa trên query parameters
    const filters: any = {};
    if (subject)
      filters.subject = { contains: subject as string, mode: "insensitive" };
    if (grade) filters.grade = parseInt(grade as string, 10);
    if (status) filters.status = status as string;

    // Tính toán skip và lấy dữ liệu
    const skip = (pageNum - 1) * pageSizeNum;
    const courses = await prisma.course.findMany({
      where: filters,
      skip,
      take: pageSizeNum,
      include: {
        tutor: {
          include: {
            profile: {
              select: {
                email: true,
                full_name: true,
                phone: true,
              },
            },
          },
        },
      },
    });

    // Đếm tổng số bản ghi
    const totalCourses = await prisma.course.count({ where: filters });

    // Trả về dữ liệu với phân trang và filter
    res.json({
      message: "Courses retrieved successfully",
      data: courses,
      pagination: {
        total: totalCourses,
        page: pageNum,
        pageSize: pageSizeNum,
        totalPages: Math.ceil(totalCourses / pageSizeNum),
      },
    });
  } catch (error: any) {
    res.status(500).json({ message: "Error retrieving courses", error });
  }
};
