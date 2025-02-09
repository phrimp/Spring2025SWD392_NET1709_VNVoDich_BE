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
exports.deleteChild = exports.updateChild = exports.createChild = exports.getChild = exports.getChildren = void 0;
const client_1 = require("@prisma/client");
const prisma = new client_1.PrismaClient();
// Lấy danh sách tất cả children của một parent
const getChildren = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { parentId } = req.params;
        if (!parentId) {
            res.status(400).json({ message: "Parent ID is required" });
            return;
        }
        const children = yield prisma.children.findMany({
            where: { parent_id: Number(parentId) },
        });
        res.json({
            message: "Children retrieved successfully",
            data: children,
        });
    }
    catch (error) {
        console.error("Error retrieving children:", error);
        res.status(500).json({ message: "Error retrieving children", error });
    }
});
exports.getChildren = getChildren;
// Lấy thông tin của một child theo ID
const getChild = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { id } = req.params;
        const { parentId } = req.params;
        if (!parentId) {
            res.status(400).json({ message: "Parent ID is required" });
            return;
        }
        if (!id || isNaN(Number(id))) {
            res.status(400).json({ message: "Invalid child ID" });
            return;
        }
        const child = yield prisma.children.findUnique({
            where: { id: Number(id) },
        });
        if (!child) {
            res.status(404).json({ message: "Child not found" });
            return;
        }
        res.json({ message: "Child retrieved successfully", data: child });
    }
    catch (error) {
        console.error("Error retrieving child:", error);
        res.status(500).json({ message: "Error retrieving child", error });
    }
});
exports.getChild = getChild;
// Tạo mới tài khoản trẻ em
const createChild = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { parentId } = req.params;
        if (!parentId) {
            res.status(400).json({ message: "Parent ID is required" });
            return;
        }
        const { full_name, age, grade_level, learning_goals, password, parent_id } = req.body;
        if (!full_name ||
            !age ||
            !grade_level ||
            !learning_goals ||
            !password ||
            !parent_id) {
            res.status(400).json({ message: "All fields are required" });
            return;
        }
        const newChild = yield prisma.children.create({
            data: {
                full_name,
                age: Number(age),
                grade_level,
                learning_goals,
                password,
                parent_id: Number(parent_id),
            },
        });
        res.json({ message: "Child account created successfully", data: newChild });
    }
    catch (error) {
        console.error("Error creating child:", error);
        res.status(500).json({ message: "Error creating child", error });
    }
});
exports.createChild = createChild;
// Cập nhật thông tin tài khoản trẻ em
const updateChild = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { parentId } = req.params;
        if (!parentId) {
            res.status(400).json({ message: "Parent ID is required" });
            return;
        }
        const { id } = req.params;
        const { full_name, age, grade_level, learning_goals, password } = req.body;
        const updatedChild = yield prisma.children.update({
            where: { id: Number(id) },
            data: { full_name, age, grade_level, learning_goals, password },
        });
        res.json({
            message: "Child account updated successfully",
            data: updatedChild,
        });
    }
    catch (error) {
        console.error("Error updating child:", error);
        res.status(500).json({ message: "Error updating child", error });
    }
});
exports.updateChild = updateChild;
// Xóa tài khoản trẻ em
const deleteChild = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { parentId } = req.params;
        if (!parentId) {
            res.status(400).json({ message: "Parent ID is required" });
            return;
        }
        const { id } = req.params;
        yield prisma.children.delete({ where: { id: Number(id) } });
        res.json({ message: "Child account deleted successfully" });
    }
    catch (error) {
        console.error("Error deleting child:", error);
        res.status(500).json({ message: "Error deleting child", error });
    }
});
exports.deleteChild = deleteChild;
