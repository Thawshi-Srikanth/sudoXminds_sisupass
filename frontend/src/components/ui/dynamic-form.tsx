"use client";
import { Controller } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { useDynamicForm } from "@/hooks/use-dynamic-form";
import type { FieldErrors } from "react-hook-form";

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

interface DynamicFormProps {
  config: FormConfig;
  onSubmit: (data: Record<string, any>) => void | Promise<void>;
}

export function DynamicForm({ config, onSubmit }: DynamicFormProps) {
  const {
    control,
    errors,
    isSubmitting,
    shouldShowField,
    getValidationRules,
    handleSubmit,
  }: {
    control: any;
    errors: FieldErrors<Record<string, any>>;
    isSubmitting: boolean;
    shouldShowField: (field: FormField) => boolean;
    getValidationRules: (field: FormField) => any;
    handleSubmit: (e?: React.BaseSyntheticEvent) => Promise<void>;
  } = useDynamicForm({
    config,
    onSubmit,
  });

  const renderField = (field: FormField) => {
    if (!shouldShowField(field)) return null;

    const fieldId = `field-${field.name}`;
    const hasError = !!errors[field.name];
    const validationRules = getValidationRules(field);

    switch (field.type) {
      case "text":
      case "email":
      case "tel":
      case "date":
        return (
          <div key={field.name} className="col-span-6 sm:col-span-3">
            <Label
              htmlFor={fieldId}
              className={cn(hasError && "text-destructive")}
            >
              {field.label}
              {field.required && (
                <span className="text-destructive ml-1">*</span>
              )}
            </Label>
            <Controller
              name={field.name}
              control={control}
              rules={validationRules}
              render={({ field: controllerField }) => (
                <Input
                  {...controllerField}
                  id={fieldId}
                  type={field.type}
                  className={cn(
                    hasError &&
                      "border-destructive focus-visible:ring-destructive"
                  )}
                />
              )}
            />
            {hasError && (
              <p className="text-sm text-destructive mt-1">
                {errors[field.name]?.message as string}
              </p>
            )}
          </div>
        );

      case "textarea":
        return (
          <div key={field.name} className="col-span-6">
            <Label
              htmlFor={fieldId}
              className={cn(hasError && "text-destructive")}
            >
              {field.label}
              {field.required && (
                <span className="text-destructive ml-1">*</span>
              )}
            </Label>
            <Controller
              name={field.name}
              control={control}
              rules={validationRules}
              render={({ field: controllerField }) => (
                <Textarea
                  {...controllerField}
                  id={fieldId}
                  className={cn(
                    hasError &&
                      "border-destructive focus-visible:ring-destructive"
                  )}
                  rows={3}
                />
              )}
            />
            {hasError && (
              <p className="text-sm text-destructive mt-1">
                {errors[field.name]?.message as string}
              </p>
            )}
          </div>
        );

      case "select":
        return (
          <div key={field.name} className="col-span-6 sm:col-span-3">
            <Label
              htmlFor={fieldId}
              className={cn(hasError && "text-destructive")}
            >
              {field.label}
              {field.required && (
                <span className="text-destructive ml-1">*</span>
              )}
            </Label>
            <Controller
              name={field.name}
              control={control}
              rules={validationRules}
              render={({ field: controllerField }) => (
                <Select
                  value={controllerField.value || ""}
                  onValueChange={controllerField.onChange}
                >
                  <SelectTrigger
                    className={cn(
                      hasError && "border-destructive focus:ring-destructive"
                    )}
                  >
                    <SelectValue
                      placeholder={`Select ${field.label.toLowerCase()}`}
                    />
                  </SelectTrigger>
                  <SelectContent>
                    {field.options?.map((option) => (
                      <SelectItem key={option} value={option}>
                        {option}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
            {hasError && (
              <p className="text-sm text-destructive mt-1">
                {errors[field.name]?.message as string}
              </p>
            )}
          </div>
        );

      case "file":
        return (
          <div key={field.name} className="col-span-6 sm:col-span-3">
            <Label
              htmlFor={fieldId}
              className={cn(hasError && "text-destructive")}
            >
              {field.label}
              {field.required && (
                <span className="text-destructive ml-1">*</span>
              )}
            </Label>
            <Controller
              name={field.name}
              control={control}
              rules={validationRules}
              render={({ field: controllerField }) => (
                <Input
                  id={fieldId}
                  type="file"
                  accept={field.accept?.join(",")}
                  onChange={(e) => {
                    const file = e.target.files?.[0];
                    controllerField.onChange(file);
                  }}
                  className={cn(
                    hasError &&
                      "border-destructive focus-visible:ring-destructive"
                  )}
                />
              )}
            />
            {hasError && (
              <p className="text-sm text-destructive mt-1">
                {errors[field.name]?.message as string}
              </p>
            )}
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <Card className="w-full max-w-4xl mx-auto">
      <CardHeader>
        <CardTitle className="text-2xl font-bold">{config.title}</CardTitle>
        <CardDescription>{config.description}</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-8">
          {config.sections.map((section, sectionIndex) => (
            <div key={sectionIndex} className="space-y-4">
              <h3 className="text-lg font-semibold text-foreground border-b pb-2">
                {section.title}
              </h3>
              <div className="grid grid-cols-6 gap-4">
                {section.fields.map(renderField)}
              </div>
            </div>
          ))}

          <div className="flex justify-end pt-6 border-t">
            <Button
              type="submit"
              size="lg"
              className="px-8"
              disabled={isSubmitting}
            >
              {isSubmitting ? "Submitting..." : "Submit Application"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
