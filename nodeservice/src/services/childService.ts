import { PrismaClient } from "@prisma/client";
import bcrypt from "bcryptjs";

const prisma = new PrismaClient();

export const getChildren = async (userId: number) => {
  return await prisma.children.findMany({
    where: {
      OR: [
        { parent_id: userId },
        {
          courseSubscriptions: {
            some: {
              course: {
                tutor_id: userId,
              },
            },
          },
        },
      ],
    },
    include: {
      profile: {
        select: { full_name: true },
      },
    },
  });
};

export const getChild = async (childId: number) => {
  return await prisma.children.findUnique({
    where: { id: childId },
    include: {
      profile: {
        select: { full_name: true },
      },
    },
  });
};

export const createChild = async (
  full_name: string,
  username: string,
  learning_goals: string,
  password: string,
  userId: number,
  date_of_birth: string
) => {
  const hashedPassword = await bcrypt.hash(password, 10);
  const newUser = await prisma.user.create({
    data: {
      full_name,
      username,
      password: hashedPassword,
      role: "Children",
    },
  });

  return await prisma.children.create({
    data: {
      id: newUser.id,
      learning_goals,
      date_of_birth,
      parent_id: userId,
    },
    include: {
      profile: {
        select: { full_name: true },
      },
    },
  });
};

export const updateChild = async (
  childId: number,
  full_name?: string,
  learning_goals?: string,
  password?: string,
  date_of_birth?: string
) => {
  const hashedPassword = password ? await bcrypt.hash(password, 10) : undefined;

  return await prisma.children.update({
    where: { id: childId },
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
        select: { full_name: true },
      },
    },
  });
};

export const deleteChild = async (childId: number) => {
  return await prisma.children.delete({
    where: { id: childId },
  });
};
