import express from "express";
import { addReviewCourse } from "../controllers/reviewCourseController";
import tutorAuth from "../middleware/tutorAuth";

const router = express.Router();

router.post("/:course_id/reviews", tutorAuth, addReviewCourse);

export default router;
