import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";
import bcrypt from "bcryptjs";

const prisma = new PrismaClient();

// Lấy danh sách tất cả children của một parent
export const getChildren = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { userId } = req.body;

    if (!userId) {
      res.status(400).json({ message: "Parent ID is required" });
      return;
    }

    console.log(userId);

    const children = await prisma.children.findMany({
      where: {
        OR: [
          { parent_id: Number(userId) },
          {
            courseSubscriptions: {
              some: {
                course: {
                  tutor_id: Number(userId),
                },
              },
            },
          },
        ],
      },
      include: {
        profile: {
          select: {
            full_name: true,
          },
        },
      },
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
    const { userId } = req.body;
    if (!userId) {
      res.status(400).json({ message: "Parent ID is required" });
      return;
    }

    if (!id || isNaN(Number(id))) {
      res.status(400).json({ message: "Invalid child ID" });
      return;
    }

    const child = await prisma.children.findUnique({
      where: { id: Number(id) },
      include: {
        profile: {
          select: {
            full_name: true,
          },
        },
      },
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
    const {
      full_name,
      username,
      learning_goals,
      password,
      userId,
      date_of_birth,
    } = req.body;

    if (
      !full_name ||
      !learning_goals ||
      !password ||
      !userId ||
      !username ||
      !date_of_birth
    ) {
      res.status(400).json({ message: "All fields are required" });
      return;
    }

    // Hash password trước khi lưu
    const saltRounds = 10;
    const hashedPassword = await bcrypt.hash(password, saltRounds);

    const newUser = await prisma.user.create({
      data: {
        full_name,
        username,
        password: hashedPassword,
        role: "Children",
      },
    });

    const newChild = await prisma.children.create({
      data: {
        id: newUser.id,
        learning_goals,
        date_of_birth,
        parent_id: Number(userId),
      },
      include: {
        profile: {
          select: {
            full_name: true,
          },
        },
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
    const { userId } = req.body;
    if (!userId) {
      res.status(400).json({ message: "Parent ID is required" });
      return;
    }
    const { id } = req.params;
    const { full_name, learning_goals, password, date_of_birth } = req.body;

    const saltRounds = 10;
    const hashedPassword = password
      ? await bcrypt.hash(password, saltRounds)
      : password;

    const updatedChild = await prisma.children.update({
      where: { id: Number(id) },
      data: {
        date_of_birth,
        learning_goals,
        profile: {
          update: {
            full_name,
            password: hashedPassword,
          },
        },
      },
      include: {
        profile: {
          select: {
            full_name: true,
          },
        },
      },
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
    const { userId } = req.body;
    if (!userId) {
      res.status(400).json({ message: "Parent ID is required" });
      return;
    }
    const { id } = req.params;

    await prisma.user.delete({ where: { id: Number(id) } });
    res.json({ message: "Child account deleted successfully" });
  } catch (error) {
    console.error("Error deleting child:", error);
    res.status(500).json({ message: "Error deleting child", error });
  }
};
