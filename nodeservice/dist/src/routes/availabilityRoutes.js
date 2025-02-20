"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = require("express");
const availabilityController_1 = require("../controllers/availabilityController");
const tutorAuth_1 = __importDefault(require("../middleware/tutorAuth"));
const router = (0, express_1.Router)();
router.get("/", tutorAuth_1.default, availabilityController_1.getTutorAvailability);
router.get("/course/:courseId", tutorAuth_1.default, availabilityController_1.getCourseAvailability);
router.put("/update", tutorAuth_1.default, availabilityController_1.updateAvailability);
exports.default = router;
