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
exports.deleteTutor = exports.createTutor = exports.getTutor = exports.getTutors = void 0;
const client_1 = require("@prisma/client");
const prisma = new client_1.PrismaClient();
//  Lấy danh sách tất cả tutors với phân trang và bộ lọc
const getTutors = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        // Lấy query từ request
        const { page = 1, pageSize = 10, isAvailable, qualifications, teachingStyle, } = req.query;
        // Chuyển đổi page và pageSize sang số nguyên
        const pageNum = parseInt(page, 10);
        const pageSizeNum = parseInt(pageSize, 10);
        // Tạo bộ lọc dựa trên query parameters
        const filters = {};
        if (isAvailable !== undefined)
            filters.is_available = isAvailable === "true";
        if (qualifications)
            filters.qualifications = { contains: qualifications };
        if (teachingStyle)
            filters.teaching_style = { contains: teachingStyle };
        // Tính toán skip và lấy dữ liệu
        const skip = (pageNum - 1) * pageSizeNum;
        const tutors = yield prisma.tutor.findMany({
            where: filters,
            skip,
            take: pageSizeNum,
        });
        // Đếm tổng số bản ghi
        const totalTutors = yield prisma.tutor.count({ where: filters });
        // Trả về dữ liệu với phân trang và filter
        res.json({
            message: "Tutors retrieved successfully",
            data: tutors,
            pagination: {
                total: totalTutors,
                page: pageNum,
                pageSize: pageSizeNum,
                totalPages: Math.ceil(totalTutors / pageSizeNum),
            },
        });
    }
    catch (error) {
        res.status(500).json({ message: "Error retrieving tutors", error });
    }
});
exports.getTutors = getTutors;
//  Lấy thông tin một tutor theo ID
const getTutor = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    const { id } = req.params;
    try {
        const tutor = yield prisma.tutor.findUnique({
            where: { id: Number(id) },
        });
        if (!tutor) {
            res.status(404).json({
                message: "Tutor not found",
            });
            return;
        }
        res.json({ message: "Tutor retrieved successfully", data: tutor });
    }
    catch (error) {
        res.status(500).json({ message: "Error retrieving tutor", error });
    }
});
exports.getTutor = getTutor;
//  Tạo một tutor mới
const createTutor = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        let { id, bio, qualifications, teaching_style, is_available, demo_video_url, image, } = req.body;
        if (!id) {
            id = Math.floor(1 + Math.random() * 9);
        }
        const newTutor = yield prisma.tutor.create({
            data: {
                id: Number(id), // ✅ ID là số
                bio,
                qualifications,
                teaching_style,
                is_available: Boolean(is_available),
                demo_video_url: demo_video_url || null,
                image: image || null,
            },
        });
        res.json({ message: "Tutor created successfully", data: newTutor });
    }
    catch (error) {
        console.error("Error creating tutor:", error);
        res.status(500).json({ message: "Error creating tutor", error });
    }
});
exports.createTutor = createTutor;
//  Xóa một tutor theo ID
const deleteTutor = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    const { id } = req.params;
    try {
        const tutor = yield prisma.tutor.delete({
            where: { id: Number(id) },
        });
        res.json({ message: "Tutor deleted successfully", data: tutor });
    }
    catch (error) {
        res.status(500).json({ message: "Error deleting tutor", error });
    }
});
exports.deleteTutor = deleteTutor;
