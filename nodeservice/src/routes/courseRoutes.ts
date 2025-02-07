import { Router } from "express";
import {
  createCourse,
  getCourse,
  getCourses,
} from "../controllers/courseController";

const router = Router();

router.get("/", getCourses);
router.get("/:id", getCourse);

router.post("/", createCourse);

export default router;
