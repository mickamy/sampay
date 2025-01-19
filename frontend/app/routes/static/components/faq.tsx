import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "~/components/ui/accordion";
import { Card, CardContent, CardTitle } from "~/components/ui/card";

export function FAQ() {
  const [openQuestion, setOpenQuestion] = useState<number | null>(null);

  const { t } = useTranslation();

  const faqItems = [
    {
      question: t("lp.faq.1.question"),
      answer: t("lp.faq.1.answer"),
    },
    {
      question: t("lp.faq.2.question"),
      answer: t("lp.faq.2.answer"),
    },
    {
      question: t("lp.faq.3.question"),
      answer: t("lp.faq.3.answer"),
    },
  ];

  return (
    <section className="my-16">
      <h2 className="text-3xl font-bold text-center mb-8">
        {t("lp.faq.title")}
      </h2>
      <Accordion
        type="single"
        collapsible
        className="max-w-2xl mx-auto space-y-4"
      >
        {faqItems.map((item, index) => (
          <Card key={item.question}>
            <AccordionItem value={item.question}>
              <AccordionTrigger className="px-4">
                <CardTitle className="text-md md:text-xl break-words whitespace-normal text-left">
                  {item.question}
                </CardTitle>
              </AccordionTrigger>
              <AccordionContent>
                <CardContent>
                  <p>{item.answer}</p>
                </CardContent>
              </AccordionContent>
            </AccordionItem>
          </Card>
        ))}
      </Accordion>
    </section>
  );
}
