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
exports.getCourses = void 0;
const client_1 = require("@prisma/client");
const prisma = new client_1.PrismaClient();
const getCourses = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        // Lấy query từ request
        const { page = 1, pageSize = 10, subject, grade, status } = req.query;
        // Chuyển đổi page và pageSize sang số nguyên
        const pageNum = parseInt(page, 10);
        const pageSizeNum = parseInt(pageSize, 10);
        // Tạo bộ lọc dựa trên query parameters
        const filters = {};
        if (subject)
            filters.subject = { contains: subject, mode: "insensitive" };
        if (grade)
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
                        user: {
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
