"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = require("express");
const multer_1 = __importDefault(require("multer"));
const courseController_1 = require("../controllers/courseController");
const router = (0, express_1.Router)();
const upload = (0, multer_1.default)({ storage: multer_1.default.memoryStorage() });
router.get("/", courseController_1.getCourses);
router.get("/:id", courseController_1.getCourse);
router.post("/", courseController_1.createCourse);
router.put("/:id", upload.single("image"), courseController_1.updateCourse);
router.put("/:courseId/add-lesson", courseController_1.addLessonToCourse);
router.put("/:courseId/update-lesson/:lessonId", courseController_1.updateLesson);
router.delete("/:courseId/delete-lesson/:lessonId", courseController_1.deleteLesson);
router.delete("/:id", courseController_1.deleteCourse);
exports.default = router;
