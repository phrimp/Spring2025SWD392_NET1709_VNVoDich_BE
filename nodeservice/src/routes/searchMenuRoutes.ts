import express from "express";
import {
  searchCourses,
  searchTutors,
  filterByPrice,
  searchByGrade,
  searchBySubject,
  filterTutorsByRating,
} from "../controllers/searchMenuController";

const router = express.Router();

router.get("/tutor", searchTutors);
router.get("/course", searchCourses);
router.get("/price", filterByPrice);
router.get("/grade", searchByGrade);
router.get("/subject", searchBySubject);
router.get("/rating", filterTutorsByRating);

export default router;

//hehee
