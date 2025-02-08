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
exports.getTutor = exports.getTutors = void 0;
const client_1 = require("@prisma/client");
const prisma = new client_1.PrismaClient();
//  Lấy danh sách tất cả tutors với phân trang và bộ lọc
const getTutors = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { page = 1, pageSize = 10, qualifications, teachingStyle, isAvailable, } = req.query;
        const pageNum = parseInt(page, 10);
        const pageSizeNum = parseInt(pageSize, 10);
        const filters = {};
        if (qualifications)
            filters.qualifications = { contains: qualifications };
        if (teachingStyle)
            filters.teaching_style = { contains: teachingStyle };
        if (isAvailable !== undefined)
            filters.is_available = isAvailable === "true";
        const skip = (pageNum - 1) * pageSizeNum;
        const tutors = yield prisma.tutor.findMany({
            where: filters,
            skip,
            take: pageSizeNum,
        });
        const totalTutors = yield prisma.tutor.count({ where: filters });
        // Tìm thông tin user tương ứng cho mỗi tutor
        const formattedTutors = yield Promise.all(tutors.map((tutor) => __awaiter(void 0, void 0, void 0, function* () {
            const user = yield prisma.user.findUnique({
                where: { id: tutor.id },
            });
            return {
                id: tutor.id,
                bio: tutor.bio,
                qualifications: tutor.qualifications,
                teaching_style: tutor.teaching_style,
                is_available: tutor.is_available,
                demo_video_url: tutor.demo_video_url || null,
                image: tutor.image || null,
                user: user
                    ? {
                        email: user.email,
                        full_name: user.full_name,
                        phone: user.phone || null,
                        google_id: user.google_id || null,
                        timezone: user.timezone,
                    }
                    : null,
            };
        })));
        res.json({
            message: "Tutors retrieved successfully",
            data: formattedTutors,
            pagination: {
                total: totalTutors,
                page: pageNum,
                pageSize: pageSizeNum,
                totalPages: Math.ceil(totalTutors / pageSizeNum),
            },
        });
    }
    catch (error) {
        console.error("Error retrieving tutors:", error);
        res.status(500).json({ message: "Error retrieving tutors", error });
    }
});
exports.getTutors = getTutors;
//  Lấy thông tin một tutor theo ID
const getTutor = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    const { id } = req.params; // ID của tutor (chính là user ID)
    try {
        const tutor = yield prisma.tutor.findUnique({
            where: { id: Number(id) },
        });
        if (!tutor) {
            res.status(404).json({ message: "Tutor not found" });
            return;
        }
        // Tìm thông tin user tương ứng
        const user = yield prisma.user.findUnique({
            where: { id: tutor.id },
        });
        const formattedTutor = {
            id: tutor.id,
            bio: tutor.bio,
            qualifications: tutor.qualifications,
            teaching_style: tutor.teaching_style,
            is_available: tutor.is_available,
            demo_video_url: tutor.demo_video_url || null,
            image: tutor.image || null,
            user: user
                ? {
                    email: user.email,
                    full_name: user.full_name,
                    phone: user.phone || null,
                    google_id: user.google_id || null,
                    timezone: user.timezone,
                }
                : null,
        };
        res.json({ message: "Tutor retrieved successfully", data: formattedTutor });
    }
    catch (error) {
        console.error("Error retrieving tutor:", error);
        res.status(500).json({ message: "Error retrieving tutor", error });
    }
});
exports.getTutor = getTutor;
//  Tạo một tutor mới
// export const createTutor = async (
//   req: Request,
//   res: Response
// ): Promise<void> => {
//   try {
//     let {
//       id,
//       bio,
//       qualifications,
//       teaching_style,
//       is_available,
//       demo_video_url,
//       image,
//     } = req.body;
//     if (!id) {
//       id = Math.floor(1 + Math.random() * 9);
//     }
//     const newTutor = await prisma.tutor.create({
//       data: {
//         id: Number(id), // ✅ ID là số
//         bio,
//         qualifications,
//         teaching_style,
//         is_available: Boolean(is_available),
//         demo_video_url: demo_video_url || null,
//         image: image || null,
//       },
//     });
//     res.json({ message: "Tutor created successfully", data: newTutor });
//   } catch (error: any) {
//     console.error("Error creating tutor:", error);
//     res.status(500).json({ message: "Error creating tutor", error });
//   }
// };
// //  Xóa một tutor theo ID
// export const deleteTutor = async (
//   req: Request,
//   res: Response
// ): Promise<void> => {
//   const { id } = req.params;
//   try {
//     const tutor = await prisma.tutor.delete({
//       where: { id: Number(id) },
//     });
//     res.json({ message: "Tutor deleted successfully", data: tutor });
//   } catch (error) {
//     res.status(500).json({ message: "Error deleting tutor", error });
//   }
// };
