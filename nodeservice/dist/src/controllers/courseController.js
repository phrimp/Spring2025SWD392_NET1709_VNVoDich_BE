"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.updateCourse = exports.deleteCourse = exports.createCourse = exports.getCourse = exports.getCourses = void 0;
const client_1 = require("@prisma/client");
const prisma = new client_1.PrismaClient();
const getCourses = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        // Lấy query từ request
        const { page = 1, pageSize = 10, subject, grade, title, status, } = req.query;
        // Chuyển đổi page và pageSize sang số nguyên
        const pageNum = parseInt(page, 10);
        const pageSizeNum = parseInt(pageSize, 10);
        // Tạo bộ lọc dựa trên query parameters
        const filters = {};
        if (subject && subject !== "all")
            filters.subject = { contains: subject };
        if (title)
            filters.title = { contains: title };
        if (grade && grade !== "all")
            filters.grade = parseInt(grade, 10);
        if (status)
            filters.status = status;
        // Tính toán skip và lấy dữ liệu
        const skip = (pageNum - 1) * pageSizeNum;
        const courses = yield prisma.course.findMany({
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
        const totalCourses = yield prisma.course.count({ where: filters });
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
    }
    catch (error) {
        res.status(500).json({ message: "Error retrieving courses", error });
    }
});
exports.getCourses = getCourses;
const getCourse = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    const { id } = req.params;
    try {
        const course = yield prisma.course.findUnique({
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
    }
    catch (error) {
        res.status(500).json({ message: "Error retrieving courses", error });
    }
});
exports.getCourse = getCourse;
const createCourse = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { tutor_id } = req.body;
        if (!tutor_id) {
            res.status(404).json({ message: "Tutor Id is required" });
            return;
        }
        const newCourse = yield prisma.course.create({
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
    }
    catch (error) {
        res.status(500).json({ message: "Error creating course", error });
    }
});
exports.createCourse = createCourse;
const deleteCourse = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    const { id } = req.params;
    // const { userId } = getAuth(req);
    try {
        const course = yield prisma.course.delete({
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
    }
    catch (error) {
        res.status(500).json({ message: "Error deleting course", error });
    }
});
exports.deleteCourse = deleteCourse;
const updateCourse = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    const { id } = req.params;
    const updateData = Object.assign({}, req.body);
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
        const course = yield prisma.course.update({
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
    }
    catch (error) {
        res.status(500).json({ message: "Error updating course", error });
    }
});
exports.updateCourse = updateCourse;
