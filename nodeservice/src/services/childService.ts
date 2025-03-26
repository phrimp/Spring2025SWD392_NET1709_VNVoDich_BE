import { PrismaClient } from "@prisma/client";
import bcrypt from "bcryptjs";
import { calculateAge } from "../lib/utils";

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
        select: { full_name: true, username: true },
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
  const parent = await prisma.parent.findUnique({
    where: { id: userId },
    select: {
      date_of_birth: true,
    },
  });

  if (!parent) {
    throw new Error("Parent not found");
  }

  if (!parent.date_of_birth) {
    throw new Error(
      "Please update your profile with your date of birth before adding a child"
    );
  }

  if (calculateAge(date_of_birth) < 6 || calculateAge(date_of_birth) > 18) {
    throw new Error("Child must be at least 6 years old and not older than 18");
  }

  if (
    calculateAge(parent.date_of_birth.toISOString()) <
    calculateAge(date_of_birth)
  ) {
    throw new Error("Child cannot be older than parent");
  }

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
  date_of_birth?: string,
  userId?: number
) => {
  const parent = await prisma.parent.findUnique({
    where: { id: userId },
    select: {
      date_of_birth: true,
    },
  });

  if (!parent) {
    throw new Error("Parent not found");
  }

  if (!parent.date_of_birth) {
    throw new Error(
      "Please update your profile with your date of birth before adding a child"
    );
  }

  if (date_of_birth) {
    if (calculateAge(date_of_birth) < 6 || calculateAge(date_of_birth) > 18) {
      throw new Error(
        "Child must be at least 6 years old and not older than 18"
      );
    }

    if (
      calculateAge(parent.date_of_birth.toISOString()) <
      calculateAge(date_of_birth)
    ) {
      throw new Error("Child cannot be older than parent");
    }
  }

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
