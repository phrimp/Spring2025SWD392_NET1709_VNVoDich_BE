import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";

const prisma = new PrismaClient();

// Lấy danh sách tất cả children của một parent
export const getChildren = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { parentId } = req.params;

    if (!parentId) {
      res.status(400).json({ message: "Parent ID is required" });
      return;
    }

    const children = await prisma.children.findMany({
      where: { parent_id: Number(parentId) },
    });

    res.json({
      message: "Children retrieved successfully",
      data: children,
    });
  } catch (error) {
    console.error("Error retrieving children:", error);
    res.status(500).json({ message: "Error retrieving children", error });
  }
};

// Lấy thông tin của một child theo ID
export const getChild = async (req: Request, res: Response): Promise<void> => {
  try {
    const { id } = req.params;
    const { parentId } = req.params;
    if (!parentId) {
      res.status(400).json({ message: "Parent ID is required" });
      return;
    }

    if (!id || isNaN(Number(id))) {
      res.status(400).json({ message: "Invalid child ID" });
      return;
    }

    const child = await prisma.children.findUnique({
      where: { id: Number(id) },
    });

    if (!child) {
      res.status(404).json({ message: "Child not found" });
      return;
    }

    res.json({ message: "Child retrieved successfully", data: child });
  } catch (error) {
    console.error("Error retrieving child:", error);
    res.status(500).json({ message: "Error retrieving child", error });
  }
};

// Tạo mới tài khoản trẻ em
export const createChild = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { parentId } = req.params;
    if (!parentId) {
      res.status(400).json({ message: "Parent ID is required" });
      return;
    }
    const { full_name, age, grade_level, learning_goals, password, parent_id } =
      req.body;

    if (
      !full_name ||
      !age ||
      !grade_level ||
      !learning_goals ||
      !password ||
      !parent_id
    ) {
      res.status(400).json({ message: "All fields are required" });
      return;
    }

    const newChild = await prisma.children.create({
      data: {
        full_name,
        age: Number(age),
        grade_level,
        learning_goals,
        password,
        parent_id: Number(parent_id),
      },
    });

    res.json({ message: "Child account created successfully", data: newChild });
  } catch (error) {
    console.error("Error creating child:", error);
    res.status(500).json({ message: "Error creating child", error });
  }
};

// Cập nhật thông tin tài khoản trẻ em
export const updateChild = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { parentId } = req.params;
    if (!parentId) {
      res.status(400).json({ message: "Parent ID is required" });
      return;
    }
    const { id } = req.params;
    const { full_name, age, grade_level, learning_goals, password } = req.body;

    const updatedChild = await prisma.children.update({
      where: { id: Number(id) },
      data: { full_name, age, grade_level, learning_goals, password },
    });

    res.json({
      message: "Child account updated successfully",
      data: updatedChild,
    });
  } catch (error) {
    console.error("Error updating child:", error);
    res.status(500).json({ message: "Error updating child", error });
  }
};

// Xóa tài khoản trẻ em
export const deleteChild = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { parentId } = req.params;
    if (!parentId) {
      res.status(400).json({ message: "Parent ID is required" });
      return;
    }
    const { id } = req.params;

    await prisma.children.delete({ where: { id: Number(id) } });
    res.json({ message: "Child account deleted successfully" });
  } catch (error) {
    console.error("Error deleting child:", error);
    res.status(500).json({ message: "Error deleting child", error });
  }
};
