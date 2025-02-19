import express from "express";
import {
  searchTutorsAndCourses,
  //   searchSchedules,
  filterByPrice,
  searchByGrade,
  searchBySubject,
  filterTutorsByRating,
} from "../controllers/searchMenuController";

const router = express.Router();

// Định nghĩa các tuyến đường API cho chức năng tìm kiếm
router.get("/search/tutor-course", searchTutorsAndCourses); // Tìm kiếm khóa học và gia sư
// router.get("/search/schedules", searchSchedules); // Tìm kiếm lịch trình có sẵn
router.get("/search/price", filterByPrice); // Lọc theo giá
router.get("/search/grade", searchByGrade); // Tìm theo cấp học
router.get("/search/subject", searchBySubject); // Tìm theo môn học
router.get("/search/rating", filterTutorsByRating); // Lọc theo đánh giá

export default router;
