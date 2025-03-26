import { PrismaClient } from "@prisma/client";

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
