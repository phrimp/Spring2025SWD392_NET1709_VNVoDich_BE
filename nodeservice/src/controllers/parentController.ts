import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";

const prisma = new PrismaClient();

export const getParents = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { page = 1, pageSize = 10 } = req.query;
    const pageNum = parseInt(page as string, 10);
    const pageSizeNum = parseInt(pageSize as string, 10);

    const skip = (pageNum - 1) * pageSizeNum;
    const parents = await prisma.parent.findMany({
      skip,
      take: pageSizeNum,
      include: {
        childrens: {
          include: {
            profile: {
              select: {
                full_name: true,
              },
            },
          },
        },
        profile: {
          select: {
            email: true,
            full_name: true,
            phone: true,
          },
        },
      },
    });

    const totalParents = await prisma.parent.count();

    res.json({
      message: "Parents retrieved successfully",
      data: parents,
      pagination: {
        total: totalParents,
        page: pageNum,
        pageSize: pageSizeNum,
        totalPages: Math.ceil(totalParents / pageSizeNum),
      },
    });
  } catch (error) {
    console.error("Error retrieving parents:", error);
    res.status(500).json({ message: "Error retrieving parents", error });
  }
};

export const getParentById = async (
  req: Request,
  res: Response
): Promise<void> => {
  const { id } = req.params;

  try {
    const parent = await prisma.parent.findUnique({
      where: { id: Number(id) },
      include: {
        childrens: {
          include: {
            profile: {
              select: {
                full_name: true,
              },
            },
          },
        },
        profile: {
          select: {
            email: true,
            full_name: true,
            phone: true,
          },
        },
      },
    });

    if (!parent) {
      res.status(404).json({ message: "Parent not found" });
      return;
    }

    res.json({ message: "Parent retrieved successfully", data: parent });
  } catch (error) {
    console.error("Error retrieving parent:", error);
    res.status(500).json({ message: "Error retrieving parent", error });
  }
};

export const updateParentProfile = async (
  req: Request,
  res: Response
): Promise<void> => {
  const { id } = req.params;
  const updateData = { ...req.body };

  try {
    const parent = await prisma.parent.update({
      where: {
        id: Number(id),
      },
      data: {
        preferred_language: updateData.preferred_language,
        notifications_enabled: updateData.notifications_enabled,
        profile: {
          update: {
            full_name: updateData.full_name,
            phone: updateData.phone,
          },
        },
      },
    });

    if (!parent) {
      res.status(404).json({ message: "Parent  profile not found" });
      return;
    }

    res.json({ message: "Parent profile updated successfully", data: parent });
  } catch (error) {
    res.status(500).json({ message: "Error updating parent profile", error });
  }
};
