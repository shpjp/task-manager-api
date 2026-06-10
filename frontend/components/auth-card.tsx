"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { ApiError } from "@/lib/api";
import { useAuth } from "@/lib/auth-context";
import { ThemeToggle } from "./theme-toggle";
import { FieldError, FullPageSpinner, Spinner } from "./ui";

interface Field {
  name: string;
  label: string;
  type: string;
  placeholder: string;
  autoComplete: string;
}

function validateClient(
  mode: "login" | "signup",
  values: Record<string, string>
): Record<string, string> {
  const errors: Record<string, string> = {};
  if (mode === "signup" && values.name.trim().length < 2) {
    errors.name = "Name must be at least 2 characters";
  }
  if (!/^\S+@\S+\.\S+$/.test(values.email ?? "")) {
    errors.email = "Enter a valid email address";
  }
  if (mode === "signup" && (values.password ?? "").length < 8) {
    errors.password = "Password must be at least 8 characters";
  }
  if (mode === "login" && !(values.password ?? "").length) {
    errors.password = "Password is required";
  }
  return errors;
}

export function AuthCard({ mode }: { mode: "login" | "signup" }) {
  const router = useRouter();
  const { user, initializing, login, signup } = useAuth();
  const [values, setValues] = useState<Record<string, string>>({
    name: "",
    email: "",
    password: "",
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [formError, setFormError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  // Already logged in (e.g. after a refresh) — go straight to tasks.
  useEffect(() => {
    if (!initializing && user) router.replace("/tasks");
  }, [initializing, user, router]);

  if (initializing || user) return <FullPageSpinner />;

  const fields: Field[] = [
    ...(mode === "signup"
      ? [{ name: "name", label: "Name", type: "text", placeholder: "Ada Lovelace", autoComplete: "name" }]
      : []),
    { name: "email", label: "Email", type: "email", placeholder: "you@example.com", autoComplete: "email" },
    {
      name: "password",
      label: "Password",
      type: "password",
      placeholder: mode === "signup" ? "At least 8 characters" : "Your password",
      autoComplete: mode === "signup" ? "new-password" : "current-password",
    },
  ];

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const clientErrors = validateClient(mode, values);
    setErrors(clientErrors);
    setFormError("");
    if (Object.keys(clientErrors).length > 0) return;

    setSubmitting(true);
    try {
      if (mode === "signup") {
        await signup(values.name.trim(), values.email.trim(), values.password);
      } else {
        await login(values.email.trim(), values.password);
      }
      router.replace("/tasks");
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.fields) setErrors(err.fields);
        else setFormError(err.message);
      } else {
        setFormError("Something went wrong. Please try again.");
      }
      setSubmitting(false);
    }
  }

  return (
    <main className="relative flex flex-1 items-center justify-center px-4 py-12">
      <div className="absolute right-4 top-4">
        <ThemeToggle />
      </div>
      <div className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold tracking-tight text-slate-800 dark:text-slate-100">
            Taskflow
          </h1>
          <p className="mt-1 text-sm text-slate-500 dark:text-slate-400">
            {mode === "signup" ? "Create your account" : "Welcome back"}
          </p>
        </div>

        <form
          onSubmit={handleSubmit}
          noValidate
          className="space-y-4 rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900"
        >
          {fields.map((field) => (
            <div key={field.name}>
              <label
                htmlFor={field.name}
                className="mb-1 block text-sm font-medium text-slate-700 dark:text-slate-300"
              >
                {field.label}
              </label>
              <input
                id={field.name}
                type={field.type}
                placeholder={field.placeholder}
                autoComplete={field.autoComplete}
                value={values[field.name]}
                onChange={(e) =>
                  setValues((v) => ({ ...v, [field.name]: e.target.value }))
                }
                className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-800 outline-none transition focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100 dark:focus:ring-indigo-900"
              />
              <FieldError message={errors[field.name]} />
            </div>
          ))}

          {formError && (
            <p className="rounded-lg bg-red-50 px-3 py-2 text-sm font-medium text-red-700 dark:bg-red-950/40 dark:text-red-300">
              {formError}
            </p>
          )}

          <button
            type="submit"
            disabled={submitting}
            className="flex w-full items-center justify-center gap-2 rounded-lg bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white transition hover:bg-indigo-700 disabled:opacity-60"
          >
            {submitting && <Spinner className="size-4 border-white/40 border-t-white" />}
            {mode === "signup" ? "Sign up" : "Log in"}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-slate-500 dark:text-slate-400">
          {mode === "signup" ? (
            <>
              Already have an account?{" "}
              <Link
                href="/login"
                className="font-medium text-indigo-600 hover:underline dark:text-indigo-400"
              >
                Log in
              </Link>
            </>
          ) : (
            <>
              New here?{" "}
              <Link
                href="/signup"
                className="font-medium text-indigo-600 hover:underline dark:text-indigo-400"
              >
                Create an account
              </Link>
            </>
          )}
        </p>
      </div>
    </main>
  );
}
