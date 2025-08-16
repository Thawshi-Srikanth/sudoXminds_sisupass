"use client";

import { useForm } from "react-hook-form";
import { useState } from "react";

interface FormField {
  name: string;
  label: string;
  type: string;
  required: boolean;
  options?: string[];
  accept?: string[];
  show_if?: {
    [key: string]: string[];
  };
}

interface FormSection {
  title: string;
  fields: FormField[];
}

interface FormConfig {
  title: string;
  description: string;
  type: string;
  sections: FormSection[];
}

interface UseDynamicFormProps {
  config: FormConfig;
  onSubmit?: (data: Record<string, any>) => void | Promise<void>;
}

export function useDynamicForm({ config, onSubmit }: UseDynamicFormProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<Record<string, any>>({
    mode: "onBlur",
    defaultValues: {},
  });

  const {
    handleSubmit,
    watch,
    control,
    formState: { errors },
  } = form;

  const watchedValues = watch() as Record<string, any>;

  const shouldShowField = (field: FormField): boolean => {
    if (!field.show_if) return true;

    for (const [dependentField, allowedValues] of Object.entries(
      field.show_if
    )) {
      const currentValue = watchedValues[dependentField];
      if (!currentValue || !allowedValues.includes(currentValue)) {
        return false;
      }
    }
    return true;
  };

  const getValidationRules = (field: FormField) => ({
    required: field.required ? `${field.label} is required` : false,
    ...(field.type === "email" && {
      pattern: {
        value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
        message: "Invalid email address",
      },
    }),
    ...(field.type === "tel" && {
      pattern: {
        value: /^[+]?[1-9][\d]{0,15}$/,
        message: "Invalid phone number",
      },
    }),
  });

  const submitForm = async (data: Record<string, any>) => {
    if (!onSubmit) return;

    setIsSubmitting(true);
    try {
      await onSubmit(data);
    } catch (error) {
      console.error("Form submission error:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return {
    form,
    control,
    errors,
    watchedValues,
    isSubmitting,
    shouldShowField,
    getValidationRules,
    handleSubmit: handleSubmit(submitForm),
  };
}
