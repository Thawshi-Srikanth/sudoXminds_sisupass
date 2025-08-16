"use client";

import { useState, useRef } from "react";
import { useForm} from "react-hook-form";
import Webcam from "react-webcam";
import { motion } from "motion/react";
import { format } from "date-fns";
import { ArrowLeft, CameraIcon, CalendarIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dropzone,
  DropzoneContent,
  DropzoneEmptyState,
} from "@/components/ui/dropzone";
import { Calendar } from "@/components/ui/calendar";
import { cn } from "@/lib/utils";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

type KYCFormValues = {
  fullName: string;
  nic: string;
  dob: Date | undefined;
  nicFront?: File;
  nicBack?: File;
};

export default function KYCForm() {
  const form = useForm<KYCFormValues>({
    defaultValues: {
      fullName: "",
      nic: "",
      dob: undefined,
    },
  });

  const [frontFile, setFrontFile] = useState<File | null>(null);
  const [backFile, setBackFile] = useState<File | null>(null);

  const [showFrontCamera, setShowFrontCamera] = useState(false);
  const [showBackCamera, setShowBackCamera] = useState(false);

  const webcamRefFront = useRef<Webcam>(null);
  const webcamRefBack = useRef<Webcam>(null);

  const videoConstraints = { facingMode: "environment" };

  const dataURLtoFile = (dataurl: string, filename: string) => {
    const arr = dataurl.split(",");
    const mime = arr[0].match(/:(.*?);/)![1];
    const bstr = atob(arr[1]);
    let n = bstr.length;
    const u8arr = new Uint8Array(n);
    while (n--) u8arr[n] = bstr.charCodeAt(n);
    return new File([u8arr], filename, { type: mime });
  };

  const captureFront = () => {
    const imageSrc = webcamRefFront.current?.getScreenshot();
    if (imageSrc) {
      const blob = dataURLtoFile(imageSrc, "nic-front.jpg");
      setFrontFile(blob);
      form.setValue("nicFront", blob);
      setShowFrontCamera(false);
    }
  };

  const captureBack = () => {
    const imageSrc = webcamRefBack.current?.getScreenshot();
    if (imageSrc) {
      const blob = dataURLtoFile(imageSrc, "nic-back.jpg");
      setBackFile(blob);
      form.setValue("nicBack", blob);
      setShowBackCamera(false);
    }
  };

  const onSubmit = (data: KYCFormValues) => {
    console.log("Submitted data:", data);
    // send to server
  };

  return (
    <div className="flex flex-1 justify-center items-center w-full">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6, delay: 0.2 }}
        className="w-full max-w-3xl p-6"
      >
        <h1 className="text-2xl font-bold mb-2">Ready to Tap & Go?</h1>
        <p className="text-base mb-6">
          Enter your details to access your pass.
        </p>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            {/* Full Name */}
            <FormField
              control={form.control}
              name="fullName"
              rules={{ required: "Full Name is required" }}
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Full Name</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter full name" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* NIC */}
            <FormField
              control={form.control}
              name="nic"
              rules={{
                required: "NIC is required",
                pattern: {
                  value: /^(\d{9}[VXvx]|\d{12})$/,
                  message: "Invalid NIC format",
                },
              }}
              render={({ field }) => (
                <FormItem>
                  <FormLabel>NIC</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter NIC" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* DOB with custom date picker */}
            <FormField
              control={form.control}
              name="dob"
              rules={{ required: "Date of Birth is required" }}
              render={({ field }) => (
                <FormItem className="flex flex-col">
                  <FormLabel>Date of Birth</FormLabel>
                  <Popover>
                    <PopoverTrigger asChild>
                      <FormControl>
                        <Button
                          variant="outline"
                          className={cn(
                            "w-full justify-start text-left font-normal",
                            !field.value && "text-muted-foreground"
                          )}
                        >
                          <CalendarIcon className="mr-2 h-4 w-4" />
                          {field.value ? (
                            format(field.value, "PPP")
                          ) : (
                            <span>Pick a date</span>
                          )}
                        </Button>
                      </FormControl>
                    </PopoverTrigger>
                    <PopoverContent className="w-auto p-0" align="start">
                      <Calendar
                        mode="single"
                        selected={field.value}
                        onSelect={field.onChange}
                        initialFocus
                      />
                    </PopoverContent>
                  </Popover>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* NIC Front */}
            <FormField
              control={form.control}
              name="nicFront"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>NIC Front Side</FormLabel>
                  <FormControl>
                    {!showFrontCamera ? (
                      <div className="grid grid-cols-5 gap-3">
                        <Dropzone
                          maxSize={1024 * 1024 * 10}
                          minSize={1024}
                          onDrop={(files) => {
                            field.onChange(files[0]);
                            setFrontFile(files[0]);
                          }}
                          className="col-span-4"
                          onError={console.error}
                          src={frontFile ? [frontFile] : undefined}
                        >
                          <DropzoneEmptyState />
                          <DropzoneContent />
                        </Dropzone>
                        <Button
                          type="button"
                          onClick={() => setShowFrontCamera(true)}
                          size="sm"
                          className="h-full"
                          variant="outline"
                        >
                          <CameraIcon width={24} height={24} />
                        </Button>
                      </div>
                    ) : (
                      <div className="grid grid-cols-5 gap-3">
                        <Webcam
                          audio={false}
                          ref={webcamRefFront}
                          screenshotFormat="image/jpeg"
                          videoConstraints={videoConstraints}
                          className="w-full h-42 object-cover rounded border col-span-4"
                        />
                        <div className="grid gap-3">
                          <Button
                            type="button"
                            variant="secondary"
                            onClick={captureFront}
                            className=" h-full"
                          >
                            <CameraIcon width={24} height={24} />
                          </Button>
                          <Button
                            type="button"
                            variant="outline"
                            onClick={() => setShowFrontCamera(false)}
                            className=" h-full"
                          >
                            <ArrowLeft width={24} height={24} />
                          </Button>
                        </div>
                      </div>
                    )}
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* NIC Back */}
            <FormField
              control={form.control}
              name="nicBack"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>NIC Back Side</FormLabel>
                  <FormControl>
                    {!showBackCamera ? (
                      <div className="grid grid-cols-5 gap-3">
                        <Dropzone
                          maxSize={1024 * 1024 * 10}
                          minSize={1024}
                          onDrop={(files) => {
                            field.onChange(files[0]);
                            setBackFile(files[0]);
                          }}
                          className="col-span-4"
                          onError={console.error}
                          src={backFile ? [backFile] : undefined}
                        >
                          <DropzoneEmptyState />
                          <DropzoneContent />
                        </Dropzone>
                        <Button
                          type="button"
                          onClick={() => setShowBackCamera(true)}
                          size="sm"
                          className="h-full"
                          variant="outline"
                        >
                          <CameraIcon width={24} height={24} />
                        </Button>
                      </div>
                    ) : (
                      <div className="grid grid-cols-5 gap-3">
                        <Webcam
                          audio={false}
                          ref={webcamRefBack}
                          screenshotFormat="image/jpeg"
                          videoConstraints={videoConstraints}
                          className="w-full h-42 object-cover rounded border col-span-4"
                        />
                        <div className="grid gap-3">
                          <Button
                            type="button"
                            variant="secondary"
                            onClick={captureBack}
                            className=" h-full"
                          >
                            <CameraIcon width={24} height={24} />
                          </Button>
                          <Button
                            type="button"
                            variant="outline"
                            onClick={() => setShowBackCamera(false)}
                            className=" h-full"
                          >
                            <ArrowLeft width={24} height={24} />
                          </Button>
                        </div>
                      </div>
                    )}
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" size="lg" className="w-full">
              Submit
            </Button>
          </form>
        </Form>
      </motion.div>
    </div>
  );
}
