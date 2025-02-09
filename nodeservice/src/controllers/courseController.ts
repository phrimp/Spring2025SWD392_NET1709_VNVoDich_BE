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
    if (grade && grade !== "all") filters.grade = parseInt(grade as string, 10);
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

export const updateCourse = async (
  req: Request,
  res: Response
): Promise<void> => {
  const { id } = req.params;
  const updateData = { ...req.body };
  // const { userId } = getAuth(req);

  try {
    if (updateData.price) {
      const price = parseInt(updateData.price);
      if (isNaN(price)) {
        res.status(400).json({
          message: "Invaid price format",
          error: "Price must be a valid number",
        });
      }
      updateData.price = price;
    }

    if (updateData.grade) {
      const grade = parseInt(updateData.grade);
      if (isNaN(grade)) {
        res.status(400).json({
          message: "Invaid price format",
          error: "Grade must be a valid number",
        });
      }
      updateData.grade = grade;
    }

    const course = await prisma.course.update({
      where: {
        id: Number(id),
        // AND: {
        //   tutor_id: {
        //     equals: tutorId
        //   }
        // }
      },
      data: updateData,
    });

    if (!course) {
      res.status(404).json({ message: "Course not found" });
      return;
    }

    res.json({ message: "Course updated successfully", data: course });
  } catch (error) {
    res.status(500).json({ message: "Error updating course", error });
  }
};

export const addLessonToCourse = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { courseId } = req.params;
    const { title, description, learning_objectives, materials_needed } =
      req.body;

    if (!courseId) {
      res.status(404).json({ message: "Course Id is required" });
      return;
    }

    const updateCourse = await prisma.course.update({
      where: {
        id: Number(courseId),
      },
      data: {
        total_lessons: {
          increment: 1,
        },
        lessons: {
          create: {
            title,
            description,
            learning_objectives,
            materials_needed,
          },
        },
      },
    });

    res.json({ message: "Lesson added successfully", data: updateCourse });
  } catch (error) {
    res.status(500).json({ message: "Error adding lesson to course", error });
  }
};

export const updateLesson = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { courseId, lessonId } = req.params;
    const { title, description, learning_objectives, materials_needed } =
      req.body;

    if (!courseId || !lessonId) {
      res.status(404).json({ message: "Course Id or Lesson Id is required" });
      return;
    }

    const updateCourse = await prisma.course.update({
      where: {
        id: Number(courseId),
      },
      data: {
        total_lessons: {
          increment: 1,
        },
        lessons: {
          update: {
            where: { id: Number(lessonId) },
            data: { title, description, learning_objectives, materials_needed },
          },
        },
      },
    });

    res.json({ message: "Lesson added successfully", data: updateCourse });
  } catch (error) {
    res.status(500).json({ message: "Error adding lesson to course", error });
  }
};

export const deleteLesson = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { courseId, lessonId } = req.params;

    if (!courseId || !lessonId) {
      res.status(404).json({ message: "Course Id or Lesson Id is required" });
      return;
    }

    const updateCourse = await prisma.course.update({
      where: {
        id: Number(courseId),
      },
      data: {
        total_lessons: {
          decrement: 1,
        },
        lessons: {
          delete: {
            id: Number(lessonId),
          },
        },
      },
    });

    res.json({ message: "Lesson deleted successfully", data: updateCourse });
  } catch (error) {
    res.status(500).json({ message: "Error deleting lesson to course", error });
  }
};
