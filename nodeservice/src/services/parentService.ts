import { PrismaClient } from "@prisma/client";
import { calculateAge } from "../lib/utils";

const prisma = new PrismaClient();

export const getParents = async (page: number, pageSize: number) => {
  const skip = (page - 1) * pageSize;
  const parents = await prisma.parent.findMany({
    skip,
    take: pageSize,
    include: {
      childrens: {
        include: {
          profile: { select: { full_name: true } },
        },
      },
      profile: {
        select: {
          email: true,
          full_name: true,
          phone: true,
          username: true,
        },
      },
    },
  });

  const totalParents = await prisma.parent.count();

  return {
    parents,
    pagination: {
      total: totalParents,
      page,
      pageSize,
      totalPages: Math.ceil(totalParents / pageSize),
    },
  };
};

export const getParentById = async (id: number) => {
  return await prisma.parent.findUnique({
    where: { id },
    include: {
      childrens: {
        include: {
          profile: { select: { full_name: true } },
        },
      },
      profile: {
        select: {
          email: true,
          full_name: true,
          phone: true,
          username: true,
        },
      },
    },
  });
};

export const updateParentProfile = async (id: number, data: any) => {
  const { date_of_birth } = data;

  if (date_of_birth) {
    const age = calculateAge(date_of_birth);

    // Kiểm tra tuổi phải trên 20
    if (age < 20) {
      throw new Error(
        "Parent must be at least 20 years old to update profile."
      );
    }
  }

  return await prisma.parent.update({
    where: { id },
    data: {
      preferred_language: data.preferred_language,
      notifications_enabled: data.notifications_enabled,
      date_of_birth: data.date_of_birth,
      profile: {
        update: {
          full_name: data.full_name,
          phone: data.phone,
        },
      },
    },
  });
};
