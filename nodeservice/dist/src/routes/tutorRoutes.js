"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = require("express");
const tutorController_1 = require("../controllers/tutorController");
const router = (0, express_1.Router)();
router.get("/", tutorController_1.getTutors);
router.get("/:id", tutorController_1.getTutor);
router.post("/", tutorController_1.createTutor);
router.delete("/:id", tutorController_1.deleteTutor);
exports.default = router;
