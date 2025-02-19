import { log } from "console";
import { NextFunction, Request, Response } from "express";
import jwt, { JwtPayload } from "jsonwebtoken";

const tutorAuth = async (
  req: Request,
  res: Response,
  next: NextFunction
): Promise<void> => {
  const authHeader = req.headers.authorization;

  // Kiểm tra header có tồn tại và có định dạng "Bearer <token>"
  if (!authHeader || !authHeader.startsWith("Bearer ")) {
    res.status(401).json({ message: "Unauthorized" });
    return;
  }

  const token = authHeader.split(" ")[1]; // Lấy token từ "Bearer <token>"

  try {
    const tokenDecode = jwt.verify(token, String(process.env.JWT_SECRET));

    if (typeof tokenDecode === "object" && "userId" in tokenDecode) {
      req.body.userId = (tokenDecode as JwtPayload).userId;
    } else {
      res.status(401).json({
        message: "Not Authorized",
      });
      return;
    }

    next();
  } catch (error: any) {
    res.json({ message: error.message });
  }
};

export default tutorAuth;
