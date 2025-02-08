import { Router } from "express";
import {
  createCourse,
  deleteCourse,
  getCourse,
  getCourses,
} from "../controllers/courseController";

const router = Router();

router.get("/", getCourses);
router.get("/:id", getCourse);

router.post("/", createCourse);

router.delete("/:id", deleteCourse);

export default router;
