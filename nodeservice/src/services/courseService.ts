import { PrismaClient } from "@prisma/client";
import { COURSE_MESSAGES } from "../message/courseMessages";

const prisma = new PrismaClient();

// Lấy danh sách khóa học với bộ lọc và phân trang
export const getCoursesService = async (
  filters: any,
  skip: number,
  pageSizeNum: number
) => {
  const courses = await prisma.course.findMany({
    where: filters,
    skip,
    take: pageSizeNum,
    include: {
      tutor: {
        include: {
          profile: {
            select: { email: true, full_name: true, phone: true },
          },
        },
      },
    },
  });

  const totalCourses = await prisma.course.count({ where: filters });

  return { courses, totalCourses };
};

// Lấy khóa học theo ID
export const getCourseByIdService = async (id: number) => {
  return prisma.course.findUnique({
    where: { id },
    include: {
      tutor: {
        include: {
          profile: {
            select: { email: true, full_name: true, phone: true },
          },
        },
      },
      lessons: true,
      courseReviews: {
        include: {
          parent: {
            select: {
              profile: { select: { full_name: true, picture: true } },
            },
          },
        },
      },
    },
  });
};

// Tạo khóa học mới
export const createCourseService = async (tutor_id: number) => {
  return prisma.course.create({
    data: {
      tutor_id,
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
};

// Cập nhật khóa học
export const updateCourseService = async (id: number, updateData: any) => {
  if (updateData.price) {
    const price = parseInt(updateData.price);

    updateData.price = price;
  }

  if (updateData.grade) {
    const grade = parseInt(updateData.grade);

    updateData.grade = grade;
  }

  const course = await prisma.course.findUnique({
    where: { id },
  });

  if (!course) {
    throw new Error(COURSE_MESSAGES.COURSE_NOT_FOUND);
  }

  if (course.total_lessons < 5) {
    throw new Error(COURSE_MESSAGES.LESSON_LESS_THAN_5);
  }

  return prisma.course.update({
    where: { id },
    data: updateData,
  });
};

// Xóa khóa học
export const deleteCourseService = async (id: number) => {
  return prisma.course.delete({ where: { id } });
};

// Thêm bài học vào khóa học
export const addLessonToCourseService = async (
  courseId: number,
  lessonData: any
) => {
  return prisma.course.update({
    where: { id: courseId },
    data: {
      total_lessons: { increment: 1 },
      lessons: { create: lessonData },
    },
  });
};

// Cập nhật bài học
export const updateLessonService = async (
  courseId: number,
  lessonId: number,
  lessonData: any
) => {
  return prisma.course.update({
    where: { id: courseId },
    data: {
      lessons: {
        update: { where: { id: lessonId }, data: lessonData },
      },
    },
  });
};

// Xóa bài học khỏi khóa học
export const deleteLessonService = async (
  courseId: number,
  lessonId: number
) => {
  return prisma.course.update({
    where: { id: courseId },
    data: {
      total_lessons: { decrement: 1 },
      lessons: { delete: { id: lessonId } },
    },
  });
};
