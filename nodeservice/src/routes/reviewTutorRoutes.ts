import express from "express";
import { getTutorReviewsForParent } from "../controllers/reviewTutorController";

const router = express.Router();

router.get("/:tutor_id", getTutorReviewsForParent);

export default router;
