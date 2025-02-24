// bookingRoutes.ts
import express from "express";
import {
  bookCourse,
  createBooking,
  getSubscriptions,
} from "../controllers/bookingController";

const router = express.Router();

// Route để phụ huynh book khóa học
router.post("/book", createBooking);

// Route để lấy danh sách đăng ký của một phụ huynh
router.get("/subscriptions/:parent_id", getSubscriptions);

export default router;
