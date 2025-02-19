"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const jsonwebtoken_1 = __importDefault(require("jsonwebtoken"));
const tutorAuth = (req, res, next) => __awaiter(void 0, void 0, void 0, function* () {
    const authHeader = req.headers.authorization;
    // Kiểm tra header có tồn tại và có định dạng "Bearer <token>"
    if (!authHeader || !authHeader.startsWith("Bearer ")) {
        res.status(401).json({ message: "Unauthorized" });
        return;
    }
    const token = authHeader.split(" ")[1]; // Lấy token từ "Bearer <token>"
    try {
        const tokenDecode = jsonwebtoken_1.default.verify(token, String(process.env.JWT_SECRET));
        if (typeof tokenDecode === "object" && "userId" in tokenDecode) {
            req.body.userId = tokenDecode.userId;
        }
        else {
            res.status(401).json({
                message: "Not Authorized",
            });
            return;
        }
        next();
    }
    catch (error) {
        res.json({ message: error.message });
    }
});
exports.default = tutorAuth;
