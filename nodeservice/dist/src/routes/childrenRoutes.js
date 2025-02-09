"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = require("express");
const childrenController_1 = require("../controllers/childrenController");
const router = (0, express_1.Router)();
router.get("/", childrenController_1.getChildren);
router.get("/:id", childrenController_1.getChild);
router.post("/", childrenController_1.createChild);
router.put("/:id", childrenController_1.updateChild);
router.delete("/:id", childrenController_1.deleteChild);
exports.default = router;
