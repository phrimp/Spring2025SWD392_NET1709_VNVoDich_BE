"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
// bookingRoutes.ts
const express_1 = __importDefault(require("express"));
const bookingController_1 = require("../controllers/bookingController");
const router = express_1.default.Router();
// Route để phụ huynh book khóa học
router.post("/book", bookingController_1.bookCourse);
// Route để lấy danh sách đăng ký của một phụ huynh
router.get("/subscriptions/:parent_id", bookingController_1.getSubscriptions);
exports.default = router;
