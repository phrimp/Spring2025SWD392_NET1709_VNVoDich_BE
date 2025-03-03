import express from "express";
import { getCourseReviewsForParent } from "../controllers/reviewCourseController";

const router = express.Router();

router.get("/:course_id", getCourseReviewsForParent);

export default router;
