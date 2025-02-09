import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";

const prisma = new PrismaClient();

export const getCourses = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    // Lấy query từ request
    const {
      page = 1,
      pageSize = 10,
      subject,
      grade,
      title,
      status,
    } = req.query;

    // Chuyển đổi page và pageSize sang số nguyên
    const pageNum = parseInt(page as string, 10);
    const pageSizeNum = parseInt(pageSize as string, 10);

    // Tạo bộ lọc dựa trên query parameters
    const filters: any = {};
    if (subject && subject !== "all")
      filters.subject = { contains: subject as string };
    if (title) filters.title = { contains: title as string };
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
        lessons: true,
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

export const getCourse = async (req: Request, res: Response): Promise<void> => {
  const { id } = req.params;

  try {
    const course = await prisma.course.findUnique({
      where: {
        id: Number(id),
      },
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
        lessons: true,
      },
    });
    if (!course) {
      res.status(404).json({
        message: "Course not found",
      });
      return;
    }
    res.json({ message: "Courses retrieved successfully", data: course });
  } catch (error) {
    res.status(500).json({ message: "Error retrieving courses", error });
  }
};

export const createCourse = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { tutor_id } = req.body;

    if (!tutor_id) {
      res.status(404).json({ message: "Tutor Id is required" });
      return;
    }

    const newCourse = await prisma.course.create({
      data: {
        tutor_id: Number(tutor_id),
        title: "Untitled Course",
        description: "",
        subject: "Uncategorized",
        grade: 0,
        total_lessons: 0,
        image: "",
        price: 0,
        status: "Draft",
      },
    });

    res.json({ message: "Course created successfully", data: newCourse });
  } catch (error) {
    res.status(500).json({ message: "Error creating course", error });
  }
};

export const deleteCourse = async (
  req: Request,
  res: Response
): Promise<void> => {
  const { id } = req.params;
  // const { userId } = getAuth(req);

  try {
    const course = await prisma.course.delete({
      where: {
        id: Number(id),
        // AND: {
        //   tutor_id: {
        //     equals: tutorId
        //   }
        // }
      },
    });

    if (!course) {
      res.status(404).json({ message: "Course not found" });
      return;
    }

    // if (course.teacherId !== userId) {
    //   res.status(403).json({ message: "Not authorized to delete this course" });
    // }

    res.json({ message: "Course deleted successfully" });
  } catch (error) {
    res.status(500).json({ message: "Error deleting course", error });
  }
};
