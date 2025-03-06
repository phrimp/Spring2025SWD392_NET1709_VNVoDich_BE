import express from "express";
import { addReviewCourse } from "../controllers/reviewCourseController";

const router = express.Router();

router.post("/:course_id/review", addReviewCourse);

export default router;
