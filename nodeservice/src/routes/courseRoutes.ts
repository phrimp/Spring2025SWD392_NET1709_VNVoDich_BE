import { Router } from "express";
import multer from "multer";
import {
  addLessonToCourse,
  createCourse,
  deleteCourse,
  deleteLesson,
  getCourse,
  getCourses,
  updateCourse,
  updateLesson,
} from "../controllers/courseController";

const router = Router();
const upload = multer({ storage: multer.memoryStorage() });

router.get("/", getCourses);
router.get("/:id", getCourse);

router.post("/", createCourse);
router.put("/:id", upload.single("image"), updateCourse);

router.put("/:courseId/add-lesson", addLessonToCourse);
router.put("/:courseId/update-lesson/:lessonId", updateLesson);
router.delete("/:courseId/delete-lesson/:lessonId", deleteLesson);

router.delete("/:id", deleteCourse);

export default router;
