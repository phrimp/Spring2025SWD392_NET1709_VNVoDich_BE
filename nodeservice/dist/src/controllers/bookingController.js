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
exports.getSubscriptions = exports.bookCourse = void 0;
// bookingController.ts
const client_1 = require("@prisma/client");
const prisma = new client_1.PrismaClient();
// Parent booking a course
const bookCourse = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { parent_id, course_id, children_ids } = req.body;
        if (!parent_id || isNaN(Number(parent_id))) {
            res.status(400).json({ message: "Invalid parent ID" });
            return;
        }
        const parent = yield prisma.user.findFirst({
            where: { id: Number(parent_id), role: "Parent" },
        });
        if (!parent) {
            res.status(403).json({ message: "Only parents can book courses." });
            return;
        }
        const course = yield prisma.course.findFirst({
            where: { id: Number(course_id), status: "Published" },
        });
        if (!course) {
            res.status(404).json({ message: "Course not found or not available." });
            return;
        }
        const tutor = yield prisma.user.findUnique({
            where: { id: course.tutor_id, role: "Tutor" },
        });
        if (!tutor) {
            res.status(404).json({ message: "Tutor not found." });
            return;
        }
        const subscriptions = yield Promise.all(children_ids.map((child_id) => __awaiter(void 0, void 0, void 0, function* () {
            return yield prisma.courseSubscription.create({
                data: {
                    status: "Active",
                    sessions_remaining: course.total_lessons,
                    course_id: Number(course_id),
                    children_id: Number(child_id),
                },
            });
        })));
        res.json({
            message: "Course booked successfully",
            tutor,
            subscriptions,
        });
    }
    catch (error) {
        console.error("Error booking course:", error);
        res.status(500).json({ message: "Error booking course", error });
    }
});
exports.bookCourse = bookCourse;
// Get subscriptions for a parent
const getSubscriptions = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { parent_id } = req.params;
        const children = yield prisma.children.findMany({
            where: { parent_id: Number(parent_id) },
            include: {
                courseSubscriptions: {
                    include: { course: true },
                },
            },
        });
        res.json({
            message: "Subscriptions retrieved successfully",
            data: children,
        });
    }
    catch (error) {
        console.error("Error retrieving subscriptions:", error);
        res.status(500).json({ message: "Error retrieving subscriptions", error });
    }
});
exports.getSubscriptions = getSubscriptions;
