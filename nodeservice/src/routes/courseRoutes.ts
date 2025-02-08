import { Router } from "express";
import multer from "multer";
import {
  createCourse,
  deleteCourse,
  getCourse,
  getCourses,
  updateCourse,
} from "../controllers/courseController";

const router = Router();
const upload = multer({ storage: multer.memoryStorage() });

router.get("/", getCourses);
router.get("/:id", getCourse);

router.post("/", createCourse);
router.put("/:id", upload.single("image"), updateCourse);

router.delete("/:id", deleteCourse);

export default router;
