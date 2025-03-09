import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";

const prisma = new PrismaClient();

// Lấy danh sách tất cả người dùng với phân trang và bộ lọc
export const getAllUsers = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const {
      page = "1",
      pageSize = "10",
      email,
      username,
      fullName,
      phone,
      role,
      status,
    } = req.query;

    const pageNum = parseInt(page as string, 10);
    const pageSizeNum = parseInt(pageSize as string, 10);

    const filters: any = {};
    if (email) filters.email = { contains: email as string };
    if (username) filters.username = { contains: username as string };
    if (fullName) filters.full_name = { contains: fullName as string };
    if (phone) filters.phone = { contains: phone as string };
    if (role) filters.role = role;
    if (status !== undefined) filters.status = status === "true";

    const skip = (pageNum - 1) * pageSizeNum;
    const users = await prisma.user.findMany({
      where: filters,
      skip,
      take: pageSizeNum,
      select: {
        email: true,
        username: true,
        full_name: true,
        phone: true,
        picture: true,
        status: true,
        role: true,
        last_login_at: true,
        created_at: true,
        updated_at: true,
        is_verified: true,
      },
    });

    
    const formattedUsers = users.map((user) => ({
      ...user,
      last_login_at: user.last_login_at?.toString(),
      created_at: user.created_at?.toString(),
      updated_at: user.updated_at?.toString(),
    }));

    const totalUsers = await prisma.user.count({ where: filters });

    res.json({
      message: "Users retrieved successfully",
      data: formattedUsers,
      pagination: {
        total: totalUsers,
        page: pageNum,
        pageSize: pageSizeNum,
        totalPages: Math.ceil(totalUsers / pageSizeNum),
      },
    });
  } catch (error) {
    console.error("Error retrieving users:", error);
    res.status(500).json({ message: "Error retrieving users", error });
  }
};

// Lấy thông tin chi tiết của một người dùng theo ID
export const getUserById = async (
  req: Request,
  res: Response
): Promise<void> => {
  const { id } = req.params;

  try {
    const user = await prisma.user.findUnique({
      where: { id: Number(id) },
      select: {
        email: true,
        username: true,
        full_name: true,
        phone: true,
        picture: true,
        status: true,
        role: true,
        last_login_at: true,
        created_at: true,
        updated_at: true,
        is_verified: true,
      },
    });

    if (!user) {
      res.status(404).json({ message: "User not found" });
      return;
    }

    res.json({
      message: "User retrieved successfully",
      data: {
        ...user,
        last_login_at: user.last_login_at?.toString(),
        created_at: user.created_at?.toString(),
        updated_at: user.updated_at?.toString(),
      },
    });
  } catch (error) {
    console.error("Error retrieving user:", error);
    res.status(500).json({ message: "Error retrieving user", error });
  }
};

// Cập nhật thông tin người dùng
export const updateUser = async (
  req: Request,
  res: Response
): Promise<void> => {
  const { id } = req.params;
  const { full_name, phone, picture, status, role, is_verified } = req.body;

  try {
    const updatedUser = await prisma.user.update({
      where: { id: Number(id) },
      data: { full_name, phone, picture, status, role, is_verified },
      select: {
        email: true,
        username: true,
        full_name: true,
        phone: true,
        picture: true,
        status: true,
        role: true,
        last_login_at: true,
        created_at: true,
        updated_at: true,
        is_verified: true,
      },
    });

    res.json({
      message: "User updated successfully",
      data: {
        ...updatedUser,
        last_login_at: updatedUser.last_login_at?.toString(),
        created_at: updatedUser.created_at?.toString(),
        updated_at: updatedUser.updated_at?.toString(),
      },
    });
  } catch (error) {
    console.error("Error updating user:", error);
    res.status(500).json({ message: "Error updating user", error });
  }
};
