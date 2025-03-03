import { Router } from "express";
import { getParents, getParentById } from "../controllers/parentController";

const router = Router();

router.get("/", getParents);
router.get("/:id", getParentById);

export default router;
